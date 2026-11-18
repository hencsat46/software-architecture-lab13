#!/usr/bin/env python3
"""Real HTTP/DB tests. Requires running compose stack; temporarily stops services."""
import json
import subprocess
import time
import urllib.error
import urllib.request
import uuid
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
BASE = {'users': 'http://localhost:8081/api/v1/users',
        'orders': 'http://localhost:8082/api/v1/orders',
        'notifications': 'http://localhost:8083/api/v1/notifications'}
checks = 0
created = []

def compose(*args, check=True):
    return subprocess.run(['docker', 'compose', *args], cwd=ROOT, check=check,
                          capture_output=True, text=True)

def request(scope, method='GET', body=None, ident='', expected=200):
    global checks
    data = json.dumps(body).encode() if body is not None else None
    req = urllib.request.Request(BASE[scope] + ('/' + ident if ident else ''),
                                 data=data, method=method,
                                 headers={'Content-Type': 'application/json'})
    try:
        with urllib.request.urlopen(req, timeout=10) as res:
            code, payload = res.status, json.load(res)
    except urllib.error.HTTPError as err:
        code, payload = err.code, json.load(err)
    assert code == expected, (scope, method, code, expected, payload)
    checks += 1
    if method == 'POST' and code == 201:
        created.append((scope, payload['data']['id']))
    return payload.get('data', payload)

def wait(scope):
    for _ in range(60):
        try:
            request(scope)
            return
        except (OSError, AssertionError):
            time.sleep(.5)
    raise RuntimeError(f'{scope} did not start')

try:
    for scope in BASE:
        wait(scope)
    user = request('users', 'POST', {'name': 'Integration', 'email': 'integration@example.com'}, expected=201)
    uid = user['id']
    request('users', ident=uid)
    request('users', 'PUT', {'name': 'Updated', 'email': 'updated@example.com'}, uid)
    request('users', 'POST', {'name': '', 'email': 'bad'}, expected=400)
    request('users', ident='not-uuid', expected=400)
    request('users', ident=str(uuid.uuid4()), expected=404)
    order_body = {'user_id': uid, 'product': 'Book', 'quantity': 2}
    before = len(request('notifications'))
    order = request('orders', 'POST', order_body, expected=201)
    assert order['notification_sent'] is True
    oid = order['id']
    request('orders', ident=oid)
    notes = request('notifications')
    assert len(notes) == before + 1
    note = next(n for n in notes if oid in n['message'])
    created.append(('notifications', note['id']))
    assert note['user_id'].lower() == uid
    request('notifications', ident=note['id'])
    request('orders', 'PUT', {**order_body, 'quantity': 3}, oid)
    request('orders', 'POST', {**order_body, 'quantity': 0}, expected=400)
    before_orders = len(request('orders'))
    request('orders', 'POST', {**order_body, 'user_id': str(uuid.uuid4())}, expected=400)
    assert len(request('orders')) == before_orders
    note2 = request('notifications', 'POST', {'user_id': uid, 'message': 'Direct stub'}, expected=201)
    request('notifications', 'PUT', {'user_id': uid, 'message': 'Updated stub'}, note2['id'])
    request('notifications', 'POST', {'user_id': uid, 'message': ''}, expected=400)
    print('PASS: CRUD, validation, HTTP order -> user -> notification scenario')

    compose('stop', 'notification-service')
    offline_order = request('orders', 'POST', order_body, expected=201)
    assert offline_order['notification_sent'] is False
    request('orders', ident=offline_order['id'])
    request('users', ident=uid)
    print('PASS: notification outage preserves order, returns 201 and notification_sent=false')
    compose('start', 'notification-service')
    wait('notifications')
    assert len(request('notifications')) == before + 2  # no automatic retry for offline order

    compose('stop', 'user-service')
    before_orders = len(request('orders'))
    request('orders', 'POST', order_body, expected=503)
    assert len(request('orders')) == before_orders
    request('orders', ident=oid)
    request('notifications')
    print('PASS: user outage returns 503, creates no order; other services stay available')
    compose('start', 'user-service')
    wait('users')
    request('users', ident=uid)
    compose('restart', 'order-service', 'notification-service')
    wait('orders')
    wait('notifications')
    request('orders', ident=oid)
    request('notifications', ident=note['id'])
    print('PASS: application restarts preserve data')

    for role, foreign in [('users', 'orders'), ('orders', 'notifications'), ('notifications', 'users')]:
        res = compose('exec', '-T', '-e', f'PGPASSWORD={role}_local', 'postgres',
                      'psql', '-h', '127.0.0.1', '-U', role, '-d', foreign, '-c', 'SELECT 1', check=False)
        assert res.returncode != 0 and 'permission denied for database' in res.stderr, res.stderr
    print('PASS: all three DB roles denied connections to another service database')

    for scope in BASE:
        request(scope, 'PUT', ({'name': 'Missing', 'email': 'm@example.com'} if scope == 'users' else
                              order_body if scope == 'orders' else {'user_id': uid, 'message': 'Missing'}),
                str(uuid.uuid4()), expected=404)
        request(scope, 'DELETE', ident=str(uuid.uuid4()), expected=404)
    for scope, ident in reversed(created):
        request(scope, 'DELETE', ident=ident)
        request(scope, ident=ident, expected=404)
    created.clear()
    print(f'PASS: {checks} HTTP status assertions; test records removed')
finally:
    compose('start', 'user-service', 'notification-service', check=False)
    for scope, ident in reversed(created):
        try:
            request(scope, 'DELETE', ident=ident)
        except Exception:
            pass
