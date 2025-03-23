document.addEventListener("DOMContentLoaded", function () {
    document.getElementById("loginForm").addEventListener("submit", async function (event) {
        event.preventDefault();

        let username = document.getElementById("username").value;
        let password = document.getElementById("password").value;

        let response = await fetch("http://localhost:441/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ username, password })
        });

        let result = await response.json();

        if (response.ok) {
            alert("Login berhasil!");
            window.location.href = "index.html"; // Redirect ke halaman validasi lisensi
        } else {
            alert("Login gagal: " + result.error);
        }
    });
});