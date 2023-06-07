export interface Application {
  id: number;
  name: string;
  script: string;
  doi: string;
  path: string;
}

export interface Output {
  step: number;
  content: string;
}
