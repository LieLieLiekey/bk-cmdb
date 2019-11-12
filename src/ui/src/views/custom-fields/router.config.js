import Meta from '@/router/meta'
import { MENU_BUSINESS_CUSTOM_FIELDS } from '@/dictionary/menu-symbol'

export default {
    name: MENU_BUSINESS_CUSTOM_FIELDS,
    path: 'custom-fields',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '自定义字段'
        }
    })
}
