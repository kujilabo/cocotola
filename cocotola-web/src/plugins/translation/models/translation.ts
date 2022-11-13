export class TranslationModel {
  id: number;
  updatedAt: string;
  text: string;
  pos: number;
  translated: string;
  lang2: string;
  provider: string;
  constructor(
    id: number,
    updatedAt: string,
    text: string,
    pos: number,
    translated: string,
    lang2: string,
    provider: string
  ) {
    this.id = id;
    this.updatedAt = updatedAt;
    this.text = text;
    this.pos = pos;
    this.translated = translated;
    this.lang2 = lang2;
    this.provider = provider;
  }
}
