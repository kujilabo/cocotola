export class WorkbookModel {
  id: number;
  version: number;
  problemType: string;
  name: string;
  questionText: string;
  subscribed: boolean;
  // problems: ProblemModel[];
  constructor(
    id: number,
    version: number,
    problemType: string,
    name: string,
    questionText: string,
    // problems: ProblemModel[],
    subscribed: boolean
  ) {
    this.id = id;
    this.version = version;
    this.problemType = problemType;
    this.name = name;
    this.questionText = questionText;
    // this.problems = problems;
    this.subscribed = subscribed;
    if (problemType == '') {
      alert('Invaid workbook');
    }
    console.log('problemType', problemType);
  }
}
