// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package footer

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func AZList() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<nav class=\"footer-az mt-12\"><h3 class=\"block mb-3 text-lg\">A-Z List</h3><p class=\"text-sm\">Search anime by alphabet name A to Z.</p><ul class=\"flex flex-wrap gap-2 mt-2\"><li><a href=\"/anime\" class=\"text-gray-400 hover:text-pink-500\">All</a></li><li><a href=\"/az-list/A\" class=\"text-gray-400 hover:text-pink-500\">A</a></li><li><a href=\"/az-list/B\" class=\"text-gray-400 hover:text-pink-500\">B</a></li><li><a href=\"/az-list/C\" class=\"text-gray-400 hover:text-pink-500\">C</a></li><li><a href=\"/az-list/A\" class=\"text-gray-400 hover:text-pink-500\">E</a></li><li><a href=\"/az-list/B\" class=\"text-gray-400 hover:text-pink-500\">F</a></li><li><a href=\"/az-list/C\" class=\"text-gray-400 hover:text-pink-500\">G</a></li><li><a href=\"/az-list/A\" class=\"text-gray-400 hover:text-pink-500\">H</a></li><li><a href=\"/az-list/B\" class=\"text-gray-400 hover:text-pink-500\">I</a></li><li><a href=\"/az-list/C\" class=\"text-gray-400 hover:text-pink-500\">J</a></li><li><a href=\"/az-list/A\" class=\"text-gray-400 hover:text-pink-500\">K</a></li><li><a href=\"/az-list/B\" class=\"text-gray-400 hover:text-pink-500\">L</a></li><li><a href=\"/az-list/C\" class=\"text-gray-400 hover:text-pink-500\">M</a></li><li><a href=\"/az-list/A\" class=\"text-gray-400 hover:text-pink-500\">N</a></li><li><a href=\"/az-list/B\" class=\"text-gray-400 hover:text-pink-500\">O</a></li><li><a href=\"/az-list/C\" class=\"text-gray-400 hover:text-pink-500\">P</a></li><li><a href=\"/az-list/A\" class=\"text-gray-400 hover:text-pink-500\">Q</a></li><li><a href=\"/az-list/B\" class=\"text-gray-400 hover:text-pink-500\">R</a></li><li><a href=\"/az-list/C\" class=\"text-gray-400 hover:text-pink-500\">S</a></li><li><a href=\"/az-list/A\" class=\"text-gray-400 hover:text-pink-500\">T</a></li><li><a href=\"/az-list/B\" class=\"text-gray-400 hover:text-pink-500\">U</a></li><li><a href=\"/az-list/C\" class=\"text-gray-400 hover:text-pink-500\">V</a></li><li><a href=\"/az-list/A\" class=\"text-gray-400 hover:text-pink-500\">W</a></li><li><a href=\"/az-list/B\" class=\"text-gray-400 hover:text-pink-500\">X</a></li><li><a href=\"/az-list/C\" class=\"text-gray-400 hover:text-pink-500\">Y</a></li><li><a href=\"/az-list/C\" class=\"text-gray-400 hover:text-pink-500\">Z</a></li></ul></nav>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}