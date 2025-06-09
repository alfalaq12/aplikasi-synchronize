document.addEventListener("DOMContentLoaded", () => {
    const spinner = document.getElementById("spinner");
    const loginForm = document.getElementById("loginForm");
  
    if (!loginForm) return;
  
    loginForm.addEventListener("submit", async (e) => {
      e.preventDefault();
  
      const username = document.getElementById("username").value;
      const password = document.getElementById("password").value;
  
      spinner.classList.remove("d-none"); // Tampilkan spinner
  
      try {
        const res = await fetch("/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ username, password }),
        });
  
        const data = await res.json();
  
        // Hide spinner sebelum alert muncul
        spinner.classList.add("d-none");
  
        if (res.ok) {
          localStorage.setItem("isLoggedIn", "true");
  
          Swal.fire({
            icon: "success",
            title: "Login berhasil!",
            showConfirmButton: false,
            timer: 1500
          }).then(() => {
            const isLicensed = localStorage.getItem("isLicensed");
            if (isLicensed === "true") {
              window.location.href = "/dashboard.html";
            } else {
              window.location.href = "/license.html";
            }
          });
  
        } else {
          Swal.fire({
            icon: "error",
            title: "Oops...",
            text: data.error || "Login gagal.",
          });
        }
  
      } catch (err) {
        console.error(err);
        spinner.classList.add("d-none");
        Swal.fire({
          icon: "error",
          title: "Oops!",
          text: "Terjadi kesalahan saat login.",
        });
      }
    });
  });  