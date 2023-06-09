export interface Application {
  id: string;
  name: string;
  script: string;
  doi: string;
  path: string;
  charged: boolean;
}

export interface Output {
  step: number;
  content: string;
}
