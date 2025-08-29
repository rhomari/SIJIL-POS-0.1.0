import { boot } from 'quasar/wrappers';
import { createI18n } from 'vue-i18n';
import messages from 'src/i18n/messages';

const i18n = createI18n({
  locale: 'en',
  fallbackLocale: 'en',
  messages,
  legacy: false
});

export default boot(({ app }) => {
  app.use(i18n);
});

export { i18n };
