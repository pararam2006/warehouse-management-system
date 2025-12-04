(function () {
    const API_BASE_URL = 'http://localhost:8080/api';
    const token = localStorage.getItem('wms_token');
    if (!token) {
        window.location.href = '/';
        return;
    }
    function authHeaders() {
        return { 'Authorization': 'Bearer ' + token, 'Content-Type': 'application/json' };
    }
    function setMsg(el, msg) { el.textContent = msg || ''; }

    const ordersBody = document.getElementById('ordersBody');
    const ordersMsg = document.getElementById('ordersMsg');

    function loadOrders() {
        setMsg(ordersMsg, '');
        fetch(API_BASE_URL + '/orders', { headers: { 'Authorization': 'Bearer ' + token } })
            .then(r => r.json().then(data => ({ ok: r.ok, data })))
            .then(({ ok, data }) => {
                if (!ok) {
                    setMsg(ordersMsg, data && data.error ? data.error : 'Ошибка загрузки заказов');
                    return;
                }
                ordersBody.innerHTML = '';
                (data || []).forEach(o => {
                    const tr = document.createElement('tr');
                    tr.innerHTML = `
                        <td>${o.id}</td>
                        <td>${o.customer}</td>
                        <td>${o.status}</td>
                        <td><button data-id="${o.id}" class="viewBtn" style="font-size:11px;">Просмотр</button></td>
                    `;
                    ordersBody.appendChild(tr);
                });
            })
            .catch(() => setMsg(ordersMsg, 'Не удалось загрузить заказы'));
    }

    ordersBody.addEventListener('click', function (e) {
        const id = e.target.getAttribute('data-id');
        if (!id) return;
        if (e.target.classList.contains('viewBtn')) {
            fetch(API_BASE_URL + '/orders/' + id, { headers: { 'Authorization': 'Bearer ' + token } })
                .then(r => r.json().then(data => ({ ok: r.ok, data })))
                .then(({ ok, data }) => {
                    if (!ok) {
                        alert(data && data.error ? data.error : 'Ошибка получения заказа');
                        return;
                    }
                    alert(JSON.stringify(data, null, 2));
                })
                .catch(() => alert('Не удалось получить заказ'));
        }
    });

    const orderForm = document.getElementById('orderForm');
    const orderFormMsg = document.getElementById('orderFormMsg');
    orderForm.addEventListener('submit', function (e) {
        e.preventDefault();
        setMsg(orderFormMsg, '');
        const customer = document.getElementById('customer').value.trim();
        const rawItems = document.getElementById('itemsJson').value.trim();
        if (!customer || !rawItems) {
            setMsg(orderFormMsg, 'Заполните поля клиента и позиций');
            return;
        }
        let items;
        try {
            items = JSON.parse(rawItems);
        } catch {
            setMsg(orderFormMsg, 'Невалидный JSON в позициях');
            return;
        }
        fetch(API_BASE_URL + '/orders', {
            method: 'POST',
            headers: authHeaders(),
            body: JSON.stringify({ customer, items })
        }).then(r => r.json().then(data => ({ ok: r.ok, data })))
            .then(({ ok, data }) => {
                if (!ok) {
                    setMsg(orderFormMsg, data && data.error ? data.error : 'Ошибка создания заказа');
                    return;
                }
                orderForm.reset();
                loadOrders();
            })
            .catch(() => setMsg(orderFormMsg, 'Не удалось создать заказ'));
    });

    const statusForm = document.getElementById('statusForm');
    const statusMsg = document.getElementById('statusMsg');
    statusForm.addEventListener('submit', function (e) {
        e.preventDefault();
        setMsg(statusMsg, '');
        const id = document.getElementById('stOrderId').value.trim();
        const status = document.getElementById('stStatus').value.trim();
        if (!id || !status) {
            setMsg(statusMsg, 'Заполните ID и статус');
            return;
        }
        fetch(API_BASE_URL + '/orders/' + id + '/status', {
            method: 'PUT',
            headers: authHeaders(),
            body: JSON.stringify({ status })
        }).then(r => r.json().then(data => ({ ok: r.ok, data })))
            .then(({ ok, data }) => {
                if (!ok) {
                    setMsg(statusMsg, data && data.error ? data.error : 'Ошибка обновления статуса');
                    return;
                }
                loadOrders();
            })
            .catch(() => setMsg(statusMsg, 'Не удалось обновить статус'));
    });

    loadOrders();
})();

