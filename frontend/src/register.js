(function () {
    const API_BASE_URL = 'http://localhost:8080/api';

    const form = document.getElementById('registerForm');
    const emailInput = document.getElementById('email');
    const passwordInput = document.getElementById('password');
    const roleInput = document.getElementById('role');
    const submitBtn = document.getElementById('submitBtn');
    const messageEl = document.getElementById('message');

    function setMessage(text, type) {
        messageEl.textContent = text || '';
        messageEl.className = '';
        if (!text) return;
        if (type === 'error') {
            messageEl.classList.add('error');
        } else if (type === 'success') {
            messageEl.classList.add('success');
        }
    }

    form.addEventListener('submit', async function (e) {
        e.preventDefault();
        setMessage('', '');

        const email = emailInput.value.trim();
        const password = passwordInput.value;
        const role = roleInput.value;

        if (!email || !password || password.length < 6) {
            setMessage('E-mail обязателен, пароль должен быть не менее 6 символов.', 'error');
            return;
        }

        submitBtn.disabled = true;

        try {
            const resp = await fetch(API_BASE_URL + '/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, password, role })
            });

            const data = await resp.json().catch(() => ({}));

            if (!resp.ok) {
                const errText = data && data.error ? data.error : 'Ошибка регистрации';
                setMessage(errText, 'error');
                return;
            }

            if (data && data.token) {
                localStorage.setItem('wms_token', data.token);
                if (data.user) {
                    localStorage.setItem('wms_user', JSON.stringify(data.user));
                }

                setMessage('Регистрация успешна, перенаправление...', 'success');

                setTimeout(function () {
                    window.location.href = '/dashboard';
                }, 800);
            } else {
                setMessage('Некорректный ответ сервера.', 'error');
            }
        } catch (err) {
            console.error(err);
            setMessage('Не удалось подключиться к серверу.', 'error');
        } finally {
            submitBtn.disabled = false;
        }
    });
})();

