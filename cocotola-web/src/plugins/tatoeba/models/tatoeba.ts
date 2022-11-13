export class TatoebaSentenceModel {
  sentenceNumber: number;
  lang2: string;
  text: string;
  author: string;
  constructor(
    sentenceNumber: number,
    lang2: string,
    text: string,
    author: string
  ) {
    this.sentenceNumber = sentenceNumber;
    this.lang2 = lang2;
    this.text = text;
    this.author = author;
  }
}

export class TatoebaSentencePairModel {
  src: TatoebaSentenceModel;
  dst: TatoebaSentenceModel;
  constructor(src: TatoebaSentenceModel, dst: TatoebaSentenceModel) {
    this.src = src;
    this.dst = dst;
  }
}
