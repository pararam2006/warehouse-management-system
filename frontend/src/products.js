(function () {
    const API_BASE_URL = 'http://localhost:8080/api';
    const token = localStorage.getItem('wms_token');
    if (!token) {
        window.location.href = '/';
        return;
    }

    const productsBody = document.getElementById('productsBody');
    const productsMessage = document.getElementById('productsMessage');
    const form = document.getElementById('productForm');
    const formMessage = document.getElementById('formMessage');

    const idInput = document.getElementById('productId');
    const skuInput = document.getElementById('sku');
    const nameInput = document.getElementById('name');
    const descInput = document.getElementById('description');
    const categoryInput = document.getElementById('categoryId');
    const supplierInput = document.getElementById('supplierId');
    const unitInput = document.getElementById('unit');

    let categories = [];
    let suppliers = [];

    document.getElementById('resetBtn').addEventListener('click', () => {
        idInput.value = '';
        form.reset();
        unitInput.value = 'pcs';
        categoryInput.value = '';
        supplierInput.value = '';
        formMessage.textContent = '';
    });

    function loadCategories() {
        fetch(API_BASE_URL + '/categories', { headers: authHeaders() })
            .then(r => r.json().then(data => ({ ok: r.ok, data })))
            .then(({ ok, data }) => {
                if (!ok) return;
                categories = data || [];
                categoryInput.innerHTML = '<option value="">— Не выбрано —</option>';
                categories.forEach(c => {
                    const opt = document.createElement('option');
                    opt.value = c.id;
                    opt.textContent = c.name;
                    categoryInput.appendChild(opt);
                });
            })
            .catch(() => {});
    }

    function loadSuppliers() {
        fetch(API_BASE_URL + '/suppliers', { headers: authHeaders() })
            .then(r => r.json().then(data => ({ ok: r.ok, data })))
            .then(({ ok, data }) => {
                if (!ok) return;
                suppliers = data || [];
                supplierInput.innerHTML = '<option value="">— Не выбрано —</option>';
                suppliers.forEach(s => {
                    const opt = document.createElement('option');
                    opt.value = s.id;
                    opt.textContent = s.name;
                    supplierInput.appendChild(opt);
                });
            })
            .catch(() => {});
    }

    function setMessage(el, msg) {
        el.textContent = msg || '';
    }

    function authHeaders() {
        return { 'Authorization': 'Bearer ' + token, 'Content-Type': 'application/json' };
    }

    function loadProducts() {
        setMessage(productsMessage, '');
        fetch(API_BASE_URL + '/products', { headers: authHeaders() })
            .then(r => r.json().then(data => ({ ok: r.ok, data })))
            .then(({ ok, data }) => {
                if (!ok) {
                    setMessage(productsMessage, data && data.error ? data.error : 'Ошибка загрузки товаров');
                    return;
                }
                productsBody.innerHTML = '';
                (data || []).forEach(p => {
                    const tr = document.createElement('tr');
                    tr.innerHTML = `
                        <td>${p.sku}</td>
                        <td>${p.name}</td>
                        <td>${p.category_id || ''}</td>
                        <td>${p.unit || ''}</td>
                        <td>
                            <button data-id="${p.id}" class="editBtn" style="font-size:11px;margin-right:4px;">Ред.</button>
                            <button data-id="${p.id}" class="delBtn" style="font-size:11px;background:#b91c1c;">X</button>
                        </td>
                    `;
                    productsBody.appendChild(tr);
                });
            })
            .catch(() => setMessage(productsMessage, 'Не удалось загрузить товары'));
    }

    productsBody.addEventListener('click', function (e) {
        const id = e.target.getAttribute('data-id');
        if (!id) return;

        if (e.target.classList.contains('editBtn')) {
            fetch(API_BASE_URL + '/products/' + id, { headers: authHeaders() })
                .then(r => r.json().then(data => ({ ok: r.ok, data })))
                .then(({ ok, data }) => {
                    if (!ok) {
                        setMessage(formMessage, data && data.error ? data.error : 'Ошибка загрузки товара');
                        return;
                    }
                    idInput.value = data.id || '';
                    skuInput.value = data.sku || '';
                    nameInput.value = data.name || '';
                    descInput.value = data.description || '';
                    categoryInput.value = data.category_id || '';
                    supplierInput.value = data.supplier_id || '';
                    unitInput.value = data.unit || 'pcs';
                })
                .catch(() => setMessage(formMessage, 'Не удалось загрузить товар'));
        } else if (e.target.classList.contains('delBtn')) {
            if (!confirm('Удалить товар?')) return;
            fetch(API_BASE_URL + '/products/' + id, {
                method: 'DELETE',
                headers: { 'Authorization': 'Bearer ' + token }
            })
                .then(r => {
                    if (!r.ok) return r.json().then(d => { throw new Error(d && d.error || 'Ошибка удаления'); });
                    loadProducts();
                })
                .catch(err => setMessage(formMessage, err.message));
        }
    });

    form.addEventListener('submit', function (e) {
        e.preventDefault();
        setMessage(formMessage, '');

        const dto = {
            sku: skuInput.value.trim(),
            name: nameInput.value.trim(),
            description: descInput.value.trim(),
            category_id: categoryInput.value.trim(),
            supplier_id: supplierInput.value.trim(),
            unit: unitInput.value.trim()
        };

        const id = idInput.value.trim();
        const method = id ? 'PUT' : 'POST';
        const url = API_BASE_URL + '/products' + (id ? '/' + id : '');

        fetch(url, {
            method,
            headers: authHeaders(),
            body: JSON.stringify(dto)
        })
            .then(r => r.json().then(data => ({ ok: r.ok, data })))
            .then(({ ok, data }) => {
                if (!ok) {
                    setMessage(formMessage, data && data.error ? data.error : 'Ошибка сохранения');
                    return;
                }
                form.reset();
                idInput.value = '';
                unitInput.value = 'pcs';
                loadProducts();
            })
            .catch(() => setMessage(formMessage, 'Не удалось сохранить товар'));
    });

    loadProducts();
    loadCategories();
    loadSuppliers();
})();

