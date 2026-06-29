import type { Article } from "./article"

export type Author = {
    id: number;
    documentId: string;
    name: string;
    avatar?: {
        url: string;
    };
    email: string;
    articles: Article[];
}