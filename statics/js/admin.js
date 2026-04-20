const form = document.getElementById("loginForm"); // Работаем с формой
const button = document.getElementById("submitBtn");

form.addEventListener("submit", (e) => {
    e.preventDefault(); // СТОП! Это отменяет перезагрузку страницы

    const userLogin = document.getElementById("username");
    const password = document.getElementById("password")

    const userData = {
        "username": userLogin.value,
        "password": password.value
    }
    SendParam(userData)

});

async function SendParam(data) {
    try {
        const response = await fetch("/admin/auth/", {
            method: "POST",
            headers: {
            "Content-Type": "application/json",
            },
            body: JSON.stringify(data)
        })
    
        const result = await response.json()
        window.location.href = result["redirect"]

    } catch(error) {
        console.log(error)
    }
}