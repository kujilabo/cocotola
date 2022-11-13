import { ProblemModel } from '@/models/problem';

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
  audioId: string;
  text: string;
  pos: string;
  lang2: string;
  translated: string;
  sentence1: EnglishWordProblemSentenceModel;
  constructor(
    id: number,
    version: number,
    updatedAt: string,
    number: number,
    problemType: string,
    audioId: string,
    text: string,
    pos: string,
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

    const sentences = p.properties['sentences'];
    const sentence = sentences[0];
    if (sentence) {
      sentence1.text = sentence['text'];
      sentence1.translated = sentence['translated'];
      sentence1.note = sentence['note'];
      // console.log(sentences);
      // console.log(sentences[0]['text']);
      // console.log(sentences[0]['translated']);
      // console.log(sentences[0]['note']);
    }

    return {
      id: p.id,
      version: p.version,
      updatedAt: p.updatedAt,
      number: p.number,
      problemType: p.problemType,
      audioId: String(p.properties['audioId']),
      text: String(p.properties['text']),
      pos: String(p.properties['pos']),
      lang2: String(p.properties['lang2']),
      translated: String(p.properties['translated']),
      sentence1,
    };
  }
}
