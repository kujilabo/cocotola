export class AudioModel {
  id: number;
  lang2: string;
  text: string;
  content: string;
  constructor(id: number, lang2: string, text: string, content: string) {
    this.id = id;
    this.lang2 = lang2;
    this.text = text;
    this.content = content;
  }
}
