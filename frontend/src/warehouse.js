import {fetchWithAuth} from "./utils.js";

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
    function setMsg(el, msg, isSuccess) {
        el.textContent = msg || '';
        el.className = '';
        if (msg && isSuccess) {
            el.classList.add('success');
        } else if (msg) {
            el.classList.add('message');
        }
    }

    const receiptForm = document.getElementById('receiptForm');
    const receiptMsg = document.getElementById('receiptMsg');
    receiptForm.addEventListener('submit', function (e) {
        e.preventDefault();
        setMsg(receiptMsg, '');
        const dto = {
            product_id: document.getElementById('rcProductId').value.trim(),
            supplier_id: document.getElementById('rcSupplierId').value.trim(),
            quantity: parseFloat(document.getElementById('rcQuantity').value),
            price: parseFloat(document.getElementById('rcPrice').value) || 0,
            expiry_date: document.getElementById('rcExpiry').value.trim()
        };
        fetchWithAuth(API_BASE_URL + '/warehouse/receipt', {
            method: 'POST',
            headers: authHeaders(),
            body: JSON.stringify(dto)
        }).then(res => {
            if (!res || !res.ok) {
                setMsg(receiptMsg, res && res.data && res.data.error ? res.data.error : 'Ошибка приёмки');
                return;
            }
            receiptForm.reset();
            setMsg(receiptMsg, 'Товар успешно принят', true);
        }).catch(err => setMsg(receiptMsg, err.message || 'Не удалось принять товар'));
    });

    const writeOffForm = document.getElementById('writeOffForm');
    const writeOffMsg = document.getElementById('writeOffMsg');
    writeOffForm.addEventListener('submit', function (e) {
        e.preventDefault();
        setMsg(writeOffMsg, '');
        const dto = {
            product_id: document.getElementById('woProductId').value.trim(),
            quantity: parseFloat(document.getElementById('woQuantity').value)
        };
        fetchWithAuth(API_BASE_URL + '/warehouse/write-off', {
            method: 'POST',
            headers: authHeaders(),
            body: JSON.stringify(dto)
        }).then(res => {
            if (!res || !res.ok) {
                setMsg(writeOffMsg, res && res.data && res.data.error ? res.data.error : 'Ошибка списания');
                return;
            }
            writeOffForm.reset();
            setMsg(writeOffMsg, 'Товар успешно списан', true);
        }).catch(err => setMsg(writeOffMsg, err.message || 'Не удалось списать товар'));
    });

    const reserveForm = document.getElementById('reserveForm');
    const reserveMsg = document.getElementById('reserveMsg');
    reserveForm.addEventListener('submit', function (e) {
        e.preventDefault();
        setMsg(reserveMsg, '');
        const dto = {
            product_id: document.getElementById('rsProductId').value.trim(),
            order_id: document.getElementById('rsOrderId').value.trim(),
            quantity: parseFloat(document.getElementById('rsQuantity').value)
        };
        fetchWithAuth(API_BASE_URL + '/warehouse/reserve', {
            method: 'POST',
            headers: authHeaders(),
            body: JSON.stringify(dto)
        }).then(res => {
            if (!res || !res.ok) {
                setMsg(reserveMsg, res && res.data && res.data.error ? res.data.error : 'Ошибка резервирования');
                return;
            }
            reserveForm.reset();
            setMsg(reserveMsg, 'Товар успешно зарезервирован', true);
        }).catch(err => setMsg(reserveMsg, err.message || 'Не удалось зарезервировать товар'));
    });

    const inventoryBody = document.getElementById('inventoryBody');
    const inventoryMsg = document.getElementById('inventoryMsg');
    function loadInventory() {
        setMsg(inventoryMsg, '');
        fetchWithAuth(API_BASE_URL + '/warehouse/inventory')
            .then(res => {
                if (!res || !res.ok) {
                    setMsg(inventoryMsg, res && res.data && res.data.error ? res.data.error : 'Ошибка загрузки остатков');
                    return;
                }
                inventoryBody.innerHTML = '';
                (res.data || []).forEach(it => {
                    const tr = document.createElement('tr');
                    tr.innerHTML = `<td>${it.product_id}</td><td>${it.quantity}</td>`;
                    inventoryBody.appendChild(tr);
                });
            })
            .catch(() => setMsg(inventoryMsg, 'Не удалось загрузить остатки'));
    }
    document.getElementById('refreshInventory').addEventListener('click', loadInventory);
    loadInventory();

    // Быстрое создание категории
    const quickCategoryForm = document.getElementById('quickCategoryForm');
    const categoryMsg = document.getElementById('categoryMsg');
    quickCategoryForm.addEventListener('submit', function(e) {
        e.preventDefault();
        setMsg(categoryMsg, '');
        const name = document.getElementById('quickCategoryName').value.trim();
        if (!name) {
            setMsg(categoryMsg, 'Введите название категории');
            return;
        }
        fetchWithAuth(API_BASE_URL + '/categories', {
            method: 'POST',
            headers: authHeaders(),
            body: JSON.stringify({ name, parent_id: '' })
        })
            .then(res => {
                if (!res || !res.ok) {
                    setMsg(categoryMsg, res && res.data && res.data.error ? res.data.error : 'Ошибка создания категории');
                    return;
                }
                document.getElementById('quickCategoryName').value = '';
                setMsg(categoryMsg, 'Категория "' + (res.data.name || name) + '" успешно создана (ID: ' + (res.data.id || '') + ')', true);
            })
            .catch(() => setMsg(categoryMsg, 'Не удалось создать категорию'));
    });

    // Быстрое создание поставщика
    const quickSupplierForm = document.getElementById('quickSupplierForm');
    const supplierMsg = document.getElementById('supplierMsg');
    quickSupplierForm.addEventListener('submit', function(e) {
        e.preventDefault();
        setMsg(supplierMsg, '');
        const name = document.getElementById('quickSupplierName').value.trim();
        if (!name) {
            setMsg(supplierMsg, 'Введите название поставщика');
            return;
        }
        const dto = {
            name: name,
            address: document.getElementById('quickSupplierAddress').value.trim(),
            phone: document.getElementById('quickSupplierPhone').value.trim(),
            email: document.getElementById('quickSupplierEmail').value.trim()
        };
        fetchWithAuth(API_BASE_URL + '/suppliers', {
            method: 'POST',
            headers: authHeaders(),
            body: JSON.stringify(dto)
        })
            .then(res => {
                if (!res || !res.ok) {
                    setMsg(supplierMsg, res && res.data && res.data.error ? res.data.error : 'Ошибка создания поставщика');
                    return;
                }
                quickSupplierForm.reset();
                setMsg(supplierMsg, 'Поставщик "' + (res.data.name || name) + '" успешно создан (ID: ' + (res.data.id || '') + ')', true);
            })
            .catch(() => setMsg(supplierMsg, 'Не удалось создать поставщика'));
    });
})();

