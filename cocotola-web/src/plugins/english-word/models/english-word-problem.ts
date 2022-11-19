import { ProblemModel, property2 } from '@/models/problem';

export const POS_ADJ = 1; // 形容詞
export const POS_ADV = 2; // 副詞
export const POS_CONJ = 3; // 接続詞
export const POS_DET = 4; // 限定詞
export const POS_MODAL = 5; // 動詞
export const POS_NOUN = 6; // 名詞
export const POS_PREP = 7; // 前置詞
export const POS_PRON = 8; // 代名詞
export const POS_VERB = 9; // 動詞
export const POS_OTHER = 99; // その他
export class EnglishWordProblemSentenceModel {
  text: string;
  translated: string;
  note: string;
  constructor(text: string, translated: string, note: string) {
    this.text = text;
    this.translated = translated;
    this.note = note;
  }
}
export class EnglishWordProblemModel {
  id: number;
  version: number;
  updatedAt: string;
  number: number;
  problemType: string;
  audioId: number;
  text: string;
  pos: number;
  lang2: string;
  translated: string;
  sentence1: EnglishWordProblemSentenceModel;
  constructor(
    id: number,
    version: number,
    updatedAt: string,
    number: number,
    problemType: string,
    audioId: number,
    text: string,
    pos: number,
    lang2: string,
    translated: string,
    sentence1: EnglishWordProblemSentenceModel
  ) {
    this.id = id;
    this.version = version;
    this.updatedAt = updatedAt;
    this.number = number;
    this.problemType = problemType;
    this.audioId = audioId;
    this.text = text;
    this.pos = pos;
    this.lang2 = lang2;
    this.translated = translated;
    this.sentence1 = sentence1;
  }

  static of(p: ProblemModel): EnglishWordProblemModel {
    const sentence1: EnglishWordProblemSentenceModel = {
      text: '',
      translated: '',
      note: '',
    };

    if (p.properties) {
      const properties = p.properties;
      const sentences = properties['sentences'] as property2[];
      if (sentences && sentences[0]) {
        const sentence = sentences[0] as { [key: string]: string };
        if (sentence) {
          sentence1.text = sentence['text'];
          sentence1.translated = String(sentence['translated']);
          sentence1.note = String(sentence['note']);
        }
      }
      // console.log(sentences);
      // console.log(sentences[0]['text']);
      // console.log(sentences[0]['translated']);
      // console.log(sentences[0]['note']);
      return {
        id: p.id,
        version: p.version,
        updatedAt: p.updatedAt,
        number: p.number,
        problemType: p.problemType,
        audioId: +(properties['audioId'] || 0),
        text: String(properties['text']),
        pos: +(properties['pos'] || 0),
        lang2: String(properties['lang2']),
        translated: String(properties['translated']),
        sentence1,
      };
    }

    return {
      id: p.id,
      version: p.version,
      updatedAt: p.updatedAt,
      number: p.number,
      problemType: p.problemType,
      audioId: 0,
      text: '',
      pos: 0,
      lang2: '',
      translated: '',
      sentence1,
    };
  }
}
