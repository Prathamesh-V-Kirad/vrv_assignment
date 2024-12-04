export interface Task {
    id: string;
    name: string;
    description: string;
    status: boolean;
    createdAt?: Date;
    updatedAt?:Date;
  }