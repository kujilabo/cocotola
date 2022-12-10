export class StatResultModel {
  date: string;
  mastered: number;
  answered: number;
  constructor(date: string, mastered: number, answered: number) {
    this.date = date;
    this.mastered = mastered;
    this.answered = answered;
    // console.log('memorized', memorized);
  }
}
export class StatHistoryModel {
  results: StatResultModel[];
  constructor(results: StatResultModel[]) {
    this.results = results;
  }
}
