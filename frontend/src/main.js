import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import ToastService from 'primevue/toastservice'
import router from './router'
import App from './App.vue'

// PrimeVue Components
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Dropdown from 'primevue/dropdown'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import Toast from 'primevue/toast'
import Tag from 'primevue/tag'
import Textarea from 'primevue/textarea'
import Card from 'primevue/card'
import ProgressSpinner from 'primevue/progressspinner'
import Divider from 'primevue/divider'
import Tooltip from 'primevue/tooltip'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import Calendar from 'primevue/calendar'
import Checkbox from 'primevue/checkbox'

// Styles
import './style.css'
import 'primevue/resources/themes/lara-light-blue/theme.css'
import 'primevue/resources/primevue.min.css'
import 'primeicons/primeicons.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(PrimeVue)
app.use(ToastService)

// Register PrimeVue components globally
app.component('Button', Button)
app.component('InputText', InputText)
app.component('InputNumber', InputNumber)
app.component('Dropdown', Dropdown)
app.component('DataTable', DataTable)
app.component('Column', Column)
app.component('Dialog', Dialog)
app.component('Toast', Toast)
app.component('Tag', Tag)
app.component('Textarea', Textarea)
app.component('Card', Card)
app.component('ProgressSpinner', ProgressSpinner)
app.component('Divider', Divider)
app.component('TabView', TabView)
app.component('TabPanel', TabPanel)
app.component('Calendar', Calendar)
app.component('Checkbox', Checkbox)

// Register directives
app.directive('tooltip', Tooltip)

// Global error handlers to surface runtime errors in the page (helpful for debugging blank pages)
window.addEventListener('error', (event) => {
	try {
		const el = document.getElementById('app')
		if (el) {
			el.innerHTML = `<div style="padding:24px;font-family:monospace;color:#b91c1c;background:#fff7f7;border:1px solid #fecaca;border-radius:6px;">JavaScript Error: ${String(event.error?.message || event.message)}<br/><small>${event.filename || ''}:${event.lineno || ''}:${event.colno || ''}</small></div>`
		}
	} catch (e) {
		console.error('Failed to render error overlay', e)
	}
	console.error('Unhandled error:', event.error || event.message)
})

window.addEventListener('unhandledrejection', (event) => {
	try {
		const el = document.getElementById('app')
		if (el) {
			el.innerHTML = `<div style="padding:24px;font-family:monospace;color:#b91c1c;background:#fff7f7;border:1px solid #fecaca;border-radius:6px;">Unhandled Promise Rejection: ${String(event.reason?.message || event.reason || '')}</div>`
		}
	} catch (e) {
		console.error('Failed to render rejection overlay', e)
	}
	console.error('Unhandled rejection:', event.reason)
})

try {
	app.mount('#app')
} catch (err) {
	try {
		const el = document.getElementById('app')
		if (el) {
			el.innerHTML = `<div style="padding:24px;font-family:monospace;color:#b91c1c;background:#fff7f7;border:1px solid #fecaca;border-radius:6px;">Mount Error: ${String(err.message)}</div>`
		}
	} catch (e) {
		console.error('Failed to render mount error', e)
	}
	console.error('Mount error:', err)
}
