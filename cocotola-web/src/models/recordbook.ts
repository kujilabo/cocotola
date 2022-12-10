export class RecordModel {
  problemId: number;
  level: number;
  resultPrev1: boolean;
  mastered: boolean;
  lastAnsweredAt: string;
  isReview: boolean;
  reviewLevel: number;
  constructor(
    problemId: number,
    level: number,
    resultPrev1: boolean,
    mastered: boolean,
    lastAnswererdAt: string
  ) {
    this.problemId = problemId;
    this.level = level;
    this.resultPrev1 = resultPrev1;
    this.mastered = mastered;
    this.lastAnsweredAt = lastAnswererdAt;
    this.isReview = false;
    this.reviewLevel = 0;
  }
}

export class RecordbookModel {
  id: number;
  records: RecordModel[];
  constructor(id: number, records: RecordModel[]) {
    this.id = id;
    this.records = records;
  }
}
