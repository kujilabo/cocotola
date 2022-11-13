import { use } from 'i18next';
import { initReactI18next } from 'react-i18next';

export const initI18n = () => {
  const f = async () => {
    await use(initReactI18next).init({
      resources: {
        en: {
          translation: {
            Delete: 'Delete',
            Edit: 'Edit',
            Pos: 'Pos',
            SignOut: 'Sign out',
            Translated: 'Translated',
            Update: 'Update',
            Word: 'Word',
          },
        },
        ja: {
          translation: {
            Delete: '削除する',
            Edit: '編集する',
            Pos: '品詞',
            SignOut: 'サインアウト',
            Translated: '翻訳',
            Update: '更新する',
            Word: '英単語',
          },
        },
      },
      lng: 'ja',
      fallbackLng: 'ja',
      interpolation: { escapeValue: false },
    });
  };
  f().catch(console.error);
};
