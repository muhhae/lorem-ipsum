/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./pkg/template/*.templ",
        "./pkg/component/*.templ",
        "./internal/views/**/*.templ"
    ],
    daisyui: {
        themes: ["nord", "sunset"],
    },
    plugins: [require("daisyui")],
}

