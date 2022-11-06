export class ProblemModel {
  id: number;
  version: number;
  updatedAt: string;
  number: number;
  problemType: string;
  properties: { [key: string]: any };
  constructor(
    id: number,
    version: number,
    updatedAt: string,
    number: number,
    problemType: string,
    properties: { [key: string]: any }
  ) {
    this.id = id;
    this.version = version;
    this.updatedAt = updatedAt;
    this.number = number;
    this.problemType = problemType;
    this.properties = properties;
  }
}

export const EnglishWordProblemTypeId = 'english_word';
export const EnglishSentenceProblemTypeId = 'english_sentence';
export const EnglishPhraseProblemTypeId = 'english_phrase';
export const TemplateProblemTypeId = 'template';

export class EnglishSentenceModel {
  number: string;
  text: string;
  translation: string;
  properties: { [key: string]: string };
  constructor(
    number: string,
    text: string,
    translation: string,
    properties: { [key: string]: string }
  ) {
    this.number = number;
    this.text = text;
    this.translation = translation;
    this.properties = properties;
  }
}
