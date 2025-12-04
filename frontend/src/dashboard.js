(function () {
    const API_BASE_URL = 'http://localhost:8080/api';

    const token = localStorage.getItem('wms_token');
    if (!token) {
        window.location.href = '/';
        return;
    }

    const userInfoEl = document.getElementById('userInfo');
    const logoutBtn = document.getElementById('logoutBtn');
    const messageEl = document.getElementById('message');

    function showError(msg) {
        messageEl.textContent = msg || '';
    }

    logoutBtn.addEventListener('click', function () {
        localStorage.removeItem('wms_token');
        localStorage.removeItem('wms_user');
        window.location.href = '/';
    });

    // Загружаем данные пользователя
    fetch(API_BASE_URL + '/auth/me', {
        headers: { 'Authorization': 'Bearer ' + token }
    })
        .then(r => r.json().then(data => ({ ok: r.ok, data, status: r.status })))
        .then(({ ok, data, status }) => {
            if (!ok) {
                showError(data && data.error ? data.error : 'Ошибка авторизации');
                if (status === 401) {
                    localStorage.removeItem('wms_token');
                    localStorage.removeItem('wms_user');
                    window.location.href = '/';
                }
                return;
            }
            if (data && data.email) {
                userInfoEl.textContent = data.email + ' (' + (data.role || 'роль не указана') + ')';
            }
        })
        .catch(() => showError('Не удалось получить данные пользователя'));

    // Простая загрузка метрик: считаем количество сущностей
    Promise.all([
        fetch(API_BASE_URL + '/products', { headers: { 'Authorization': 'Bearer ' + token } }),
        fetch(API_BASE_URL + '/categories', { headers: { 'Authorization': 'Bearer ' + token } }),
        fetch(API_BASE_URL + '/suppliers', { headers: { 'Authorization': 'Bearer ' + token } }),
        fetch(API_BASE_URL + '/orders', { headers: { 'Authorization': 'Bearer ' + token } })
    ]).then(async ([p, c, s, o]) => {
        const [pd, cd, sd, od] = await Promise.all([
            p.json().catch(() => []),
            c.json().catch(() => []),
            s.json().catch(() => []),
            o.json().catch(() => [])
        ]);
        document.getElementById('metricProducts').textContent = Array.isArray(pd) ? pd.length : '—';
        document.getElementById('metricCategories').textContent = Array.isArray(cd) ? cd.length : '—';
        document.getElementById('metricSuppliers').textContent = Array.isArray(sd) ? sd.length : '—';
        document.getElementById('metricOrders').textContent = Array.isArray(od) ? od.length : '—';
    }).catch(() => {
        showError('Не удалось загрузить метрики');
    });
})();

