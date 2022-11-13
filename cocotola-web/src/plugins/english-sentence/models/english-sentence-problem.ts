import { ProblemModel } from '@/models/problem';

export class EnglishSentenceProblemModel {
  id: number;
  version: number;
  updatedAt: string;
  number: number;
  problemType: string;
  provider: string;
  audioId: string;
  text: string;
  lang2: string;
  translated: string;
  // note: string;
  constructor(
    id: number,
    version: number,
    updatedAt: string,
    number: number,
    problemType: string,
    provider: string,
    audioId: string,
    text: string,
    lang2: string,
    translated: string
    // note: string
  ) {
    this.id = id;
    this.version = version;
    this.updatedAt = updatedAt;
    this.number = number;
    this.problemType = problemType;
    this.provider = provider;
    this.audioId = audioId;
    this.text = text;
    this.lang2 = lang2;
    this.translated = translated;
    // this.note = note;
  }

  static of(p: ProblemModel): EnglishSentenceProblemModel {
    return {
      id: p.id,
      version: p.version,
      updatedAt: p.updatedAt,
      number: p.number,
      problemType: p.problemType,
      provider: String(p.properties['provider']),
      audioId: String(p.properties['audioId']),
      text: String(p.properties['text']),
      lang2: String(p.properties['lang2']),
      translated: String(p.properties['translated']),
      // note: '' + p.properties['note'],
    };
  }
}
