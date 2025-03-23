document.addEventListener("DOMContentLoaded", function () {
    if (!localStorage.getItem("isLoggedIn")) {
        window.location.href = "login.html"; // Pastikan user sudah login
        return;
    }

    let licenseForm = document.getElementById("licenseForm");

    if (!licenseForm) {
        console.error("Error: Formulir lisensi tidak ditemukan!");
        return;
    }

    licenseForm.addEventListener("submit", async function (event) {
        event.preventDefault();

        let licenseKey = document.getElementById("licenseKey").value;
        if (!licenseKey) {
            alert("Silakan masukkan kunci lisensi!");
            return;
        }

        try {
            let response = await fetch("/validate-license", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ key: licenseKey })
            });

            let result = await response.text();

            if (response.ok) {
                localStorage.setItem("hasLicense", "true"); // Simpan lisensi valid
                alert("Lisensi valid! Mengalihkan ke dashboard...");
                window.location.href = "dashboard.html";
            } else {
                alert("Error: " + result);
            }
        } catch (error) {
            console.error("Gagal memvalidasi lisensi:", error);
            alert("Terjadi kesalahan saat menghubungi server.");
        }
    });
});
