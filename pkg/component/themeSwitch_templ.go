// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.501
package component

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func ThemeToggleScript() templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_ThemeToggleScript_8d8b`,
		Function: `function __templ_ThemeToggleScript_8d8b(){document.addEventListener('DOMContentLoaded', (event) => {
		var themeController = document.querySelector('.my-theme');
		if (!themeController) return;
		dark = themeController.getAttribute('dark') ? themeController.getAttribute('dark') : 'dark'
		light = themeController.getAttribute('light') ? themeController.getAttribute('light') : 'light'
		var theme = localStorage.getItem('theme');
		if (theme) {
			document.documentElement.setAttribute('data-theme', theme)
			if (theme == "dark") {
				themeController.checked = true;
				document.documentElement.setAttribute('data-theme', dark)
			} else {
				themeController.checked = false;
				document.documentElement.setAttribute('data-theme', light)
			}
		} else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
			themeController.checked = true;
			document.documentElement.setAttribute('data-theme', dark)
			localStorage.setItem('theme', "dark")
		} else {
			document.documentElement.setAttribute('data-theme', light)
			localStorage.setItem('theme', "light")
		}
		themeController.addEventListener('change', function (e) {
			if (e.currentTarget.checked) {
				document.documentElement.setAttribute('data-theme', dark)
				localStorage.setItem('theme', "dark")
			} else {
				document.documentElement.setAttribute('data-theme', light)
				localStorage.setItem('theme', "light")
			}
		})
	});}`,
		Call:       templ.SafeScript(`__templ_ThemeToggleScript_8d8b`),
		CallInline: templ.SafeScriptInline(`__templ_ThemeToggleScript_8d8b`),
	}
}

func ThemeSwitch() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<label class=\"cursor-pointer grid place-items-center\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = ThemeToggleScript().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<input type=\"checkbox\" light=\"nord\" dark=\"sunset\" class=\"toggle my-theme bg-base-content row-start-1 col-start-1 col-span-2\"> <svg class=\"col-start-1 row-start-1 stroke-base-100 fill-base-100\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"currentColor\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><circle cx=\"12\" cy=\"12\" r=\"5\"></circle><path d=\"M12 1v2M12 21v2M4.2 4.2l1.4 1.4M18.4 18.4l1.4 1.4M1 12h2M21 12h2M4.2 19.8l1.4-1.4M18.4 5.6l1.4-1.4\"></path></svg> <svg class=\"col-start-2 row-start-1 stroke-base-100 fill-base-100\" xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 24 24\" fill=\"currentColor\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z\"></path></svg></label>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}