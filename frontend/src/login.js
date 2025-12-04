(function () {
    const API_BASE_URL = 'http://localhost:8080/api';

    const form = document.getElementById('loginForm');
    const emailInput = document.getElementById('email');
    const passwordInput = document.getElementById('password');
    const submitBtn = document.getElementById('submitBtn');
    const messageEl = document.getElementById('message');

    function setMessage(el, text, type) {
        el.textContent = text || '';
        el.className = '';
        if (!text) return;
        if (type === 'error') {
            el.classList.add('error');
        } else if (type === 'success') {
            el.classList.add('success');
        }
    }

    form.addEventListener('submit', async function (e) {
        e.preventDefault();
        setMessage(messageEl, '', '');

        const email = emailInput.value.trim();
        const password = passwordInput.value;

        if (!email || !password) {
            setMessage(messageEl, 'Введите e-mail и пароль.', 'error');
            return;
        }

        submitBtn.disabled = true;

        try {
            const resp = await fetch(API_BASE_URL + '/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, password })
            });

            const data = await resp.json().catch(() => ({}));

            if (!resp.ok) {
                const errText = data && data.error ? data.error : 'Ошибка входа';
                setMessage(messageEl, errText, 'error');
                return;
            }

            if (data && data.token) {
                localStorage.setItem('wms_token', data.token);
                if (data.user) {
                    localStorage.setItem('wms_user', JSON.stringify(data.user));
                }

                setMessage(messageEl, 'Вход выполнен, перенаправление...', 'success');

                setTimeout(() => {
                    window.location.href = '/dashboard';
                }, 600);
            } else {
                setMessage(messageEl, 'Некорректный ответ сервера.', 'error');
            }
        } catch (err) {
            console.error(err);
            setMessage(messageEl, 'Не удалось подключиться к серверу.', 'error');
        } finally {
            submitBtn.disabled = false;
        }
    });
})();

