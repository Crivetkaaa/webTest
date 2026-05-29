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
        const href = result["redirect"]
        const resBool = result["success"]

        if (!resBool) {
            const errPage = document.getElementById("errMessage")
            errPage.style.color = "red"
            errPage.style.display = "flex"
            errPage.style.justifyContent = "center"
            errPage.style.alignItems = "center"
            errPage.textContent = result["message"]
        } else {
            window.location.href = href
        }

    } catch(error) {
        console.log(error)
    }
}