import { ProblemModel, propertyObject } from '@/models/problem';
export class EnglishSentenceProblemModel {
  id: number;
  version: number;
  updatedAt: string;
  number: number;
  problemType: string;
  provider: string;
  audioId: number;
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
    audioId: number,
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
    if (p.properties) {
      const properties = p.properties as propertyObject;

      return {
        id: p.id,
        version: p.version,
        updatedAt: p.updatedAt,
        number: p.number,
        problemType: p.problemType,
        provider: String(properties['provider']),
        audioId: +(properties['audioId'] || 0),
        text: String(properties['text']),
        lang2: String(properties['lang2']),
        translated: String(properties['translated']),
        // note: '' + p.properties['note'],
      };
    }

    return {
      id: p.id,
      version: p.version,
      updatedAt: p.updatedAt,
      number: p.number,
      problemType: p.problemType,
      provider: '',
      audioId: 0,
      text: '',
      lang2: '',
      translated: '',
      // note: '' + p.properties['note'],
    };
  }
}
