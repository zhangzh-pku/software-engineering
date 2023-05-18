export interface Application {
  id: number;
  name: string;
  script: string;
  doi: string;
}

export interface Output {
  step: number;
  content: string;
}
