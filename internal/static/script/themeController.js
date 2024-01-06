document.addEventListener('DOMContentLoaded', (event) => {
    var themeController = document.querySelector('.my-theme');
    if (!themeController) return;
    dark = themeController.getAttribute('dark') ? themeController.getAttribute('dark') : 'dark'
    light = themeController.getAttribute('light') ? themeController.getAttribute('light') : 'light'
    var theme = localStorage.getItem('theme');
    if (theme) {
        document.documentElement.setAttribute('data-theme', theme)
        if (theme == dark) {
            themeController.checked = true;
        }
    } else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        themeController.checked = true;
    }
    themeController.addEventListener('change', function (e) {
        if (e.currentTarget.checked) {
            document.documentElement.setAttribute('data-theme', dark)
            localStorage.setItem('theme', dark)
        } else {
            document.documentElement.setAttribute('data-theme', light)
            localStorage.setItem('theme', light)
        }
    })
});