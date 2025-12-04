import {fetchWithAuth} from "./utils.js";

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
    fetchWithAuth(API_BASE_URL + '/auth/me')
        .then(res => {
            if (!res || !res.ok) {
                showError(res && res.data && res.data.error ? res.data.error : 'Ошибка авторизации');
                if (res && res.data && res.data.error && res.data.error.includes('401')) {
                    localStorage.removeItem('wms_token');
                    localStorage.removeItem('wms_user');
                    window.location.href = '/';
                }
                return;
            }
            if (res.data && res.data.email) {
                userInfoEl.textContent = res.data.email + ' (' + (res.data.role || 'роль не указана') + ')';
            }
        })
        .catch(() => showError('Не удалось получить данные пользователя'));

    // Простая загрузка метрик: считаем количество сущностей
    Promise.all([
        fetchWithAuth(API_BASE_URL + '/products'),
        fetchWithAuth(API_BASE_URL + '/categories'),
        fetchWithAuth(API_BASE_URL + '/suppliers'),
        fetchWithAuth(API_BASE_URL + '/orders')
    ]).then(([p, c, s, o]) => {
        const pd = (p && p.ok && Array.isArray(p.data)) ? p.data : [];
        const cd = (c && c.ok && Array.isArray(c.data)) ? c.data : [];
        const sd = (s && s.ok && Array.isArray(s.data)) ? s.data : [];
        const od = (o && o.ok && Array.isArray(o.data)) ? o.data : [];
        document.getElementById('metricProducts').textContent = pd.length;
        document.getElementById('metricCategories').textContent = cd.length;
        document.getElementById('metricSuppliers').textContent = sd.length;
        document.getElementById('metricOrders').textContent = od.length;
    }).catch(() => {
        showError('Не удалось загрузить метрики');
    });
})();

