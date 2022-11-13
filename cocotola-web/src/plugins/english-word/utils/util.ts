import {
  POS_ADJ,
  POS_ADV,
  POS_CONJ,
  POS_DET,
  POS_MODAL,
  POS_NOUN,
  POS_PREP,
  POS_PRON,
  POS_VERB,
} from '../models/english-word-problem';

export const toDsiplayText = (pos: number): string => {
  switch (pos) {
    case POS_ADJ:
      return '形容詞';
    case POS_ADV:
      return '副詞';
    case POS_CONJ:
      return '接続詞';
    case POS_DET:
      return '限定詞';
    case POS_MODAL:
      return '動詞';
    case POS_NOUN:
      return '名詞';
    case POS_PREP:
      return '前置詞';
    case POS_PRON:
      return '代名詞';
    case POS_VERB:
      return '動詞';
    default:
      return '';
  }
};
