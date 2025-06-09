document.addEventListener("DOMContentLoaded", () => {
  const spinner = document.getElementById("spinner");
  const licenseForm = document.getElementById("licenseForm");

  if (!licenseForm) return;

  licenseForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const licenseKey = document.getElementById("licenseKey").value;

    // Tampilkan spinner saat mulai validasi
    spinner.classList.remove("d-none");

    try {
      const res = await fetch("/validate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ key: licenseKey }),
      });
      
      const contentType = res.headers.get("content-type");
      if (!contentType || !contentType.includes("application/json")) {
        throw new Error("Response bukan JSON");
      }
      
      const data = await res.json();      

      if (res.ok) {
        localStorage.setItem("isLicensed", "true");

        Swal.fire({
          icon: "success",
          title: "Lisensi valid!",
          showConfirmButton: false,
          timer: 1500,
        }).then(() => {
          spinner.classList.add("d-none"); // Sembunyikan spinner
          window.location.href = "/dashboard.html";
        });
      } else {
        spinner.classList.add("d-none"); // Sembunyikan spinner
        Swal.fire({
          icon: "error",
          title: "Lisensi tidak valid",
          text: data.error || "Silakan cek kembali kode lisensimu.",
        });
      }
    } catch (err) {
      console.error(err);
      spinner.classList.add("d-none"); // Sembunyikan spinner
      Swal.fire({
        icon: "error",
        title: "Oops!",
        text: "Terjadi kesalahan saat validasi lisensi.",
      });
    }
  });
});