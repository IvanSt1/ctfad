import {ToastSeverity} from 'primevue/api';
import {app} from '@/main';

const lifeTime = 2000;

export function toastInfo(body, title = "Info") {
  app.config.globalProperties.$toast.add({severity: ToastSeverity.INFO, summary: title, detail: body, life: lifeTime});
}

export function toastSuccess(body, title = "Succes") {
  app.config.globalProperties.$toast.add({severity: ToastSeverity.SUCCESS, summary: title, detail: body, life: lifeTime});
}

export function toastError(body, title = "Error") {
  app.config.globalProperties.$toast.add({severity: ToastSeverity.ERROR, summary: title, detail: body, life: lifeTime});
}
