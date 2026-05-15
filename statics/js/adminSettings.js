// Функция показа уведомлений (Toast)
function showToast(message, isSuccess = true) {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `px-4 py-3 rounded-md shadow-lg text-white font-medium transform transition-all duration-300 translate-y-5 opacity-0 ${
        isSuccess ? 'bg-green-600' : 'bg-red-600'
    }`;
    toast.textContent = message;
    container.appendChild(toast);
    
    // Анимация появления
    setTimeout(() => {
        toast.classList.remove('translate-y-5', 'opacity-0');
    }, 10);

    // Автоудаление через 4 секунды
    setTimeout(() => {
        toast.classList.add('opacity-0');
        setTimeout(() => toast.remove(), 300);
    }, 4000);
}

// ОБРАБОТКА: Загрузка документов
document.getElementById('document-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const fileOffer = this.querySelector('input[name="offer"]').files[0];
    const filePrivacy = this.querySelector('input[name="privacy"]').files[0];

    if (!fileOffer && !filePrivacy) {
        showToast('Выберите хотя бы один файл для обновления', false);
        return;
    }

    const formData = new FormData(this);
    const submitBtn = this.querySelector('button[type="submit"]');
    
    try {
        submitBtn.disabled = true;
        submitBtn.textContent = 'Отправка...';

        const response = await fetch(this.action, {
            method: 'POST',
            body: formData // FormData автоматически выставляет нужный Boundary и Content-Type
        });

        if (response.ok) {
            showToast('Документы успешно обновлены');
            this.reset();
        } else {
            const errText = await response.text();
            showToast(`Ошибка сервера: ${errText || response.statusText}`, false);
        }
    } catch (error) {
        showToast('Ошибка сети. Не удалось отправить файлы', false);
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Обновить документы';
    }
});

// ОБРАБОТКА: Изменение пароля
document.getElementById('password-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const newPassword = this.querySelector('input[name="new_password"]').value;
    const confirmPassword = this.querySelector('input[name="confirm_password"]').value;

    if (newPassword !== confirmPassword) {
        showToast('Новые пароли не совпадают!', false);
        return;
    }

    // Используем URLSearchParams для стандартной отправки x-www-form-urlencoded
    const formData = new URLSearchParams(new FormData(this));
    const submitBtn = this.querySelector('button[type="submit"]');

    try {
        submitBtn.disabled = true;
        submitBtn.textContent = 'Сохранение...';

        const response = await fetch(this.action, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            },
            body: formData
        });

        if (response.ok) {
            showToast('Пароль успешно изменен');
            this.reset();
        } else {
            const errText = await response.text();
            showToast(`Ошибка: ${errText || response.statusText}`, false);
        }
    } catch (error) {
        showToast('Ошибка сети. Не удалось изменить пароль', false);
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Изменить пароль';
    }
});
