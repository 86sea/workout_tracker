function handleAuthForm(endpoint) {
    const form = document.getElementById('auth-form');
    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = new FormData(form);
        const data = {};
        formData.forEach((value, key) => {
            data[key] = value;
        });

        try {
            const response = await fetch(endpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
            });

            const result = await response.json();
            if (result.user && result.auth) {
                window.location.href = "/index.html";
            } else {
                // Display error to user
            }
        } catch (error) {
            // Handle network or other errors
        }
    });
}
