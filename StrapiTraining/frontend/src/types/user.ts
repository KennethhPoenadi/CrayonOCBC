export type User = {
  id: number;
  documentId?: string;

  username: string;
  email: string;

  provider?: string;
  confirmed: boolean;
  blocked: boolean;

  role?: {
    id: number;
    name: string;
    description?: string;
    type?: string;
  };

  createdAt?: string;
  updatedAt?: string;
};