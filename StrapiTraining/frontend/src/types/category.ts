import type { Article } from "./article"


export type Category = {
    id: number;
    documentId: string;
    slug: string;
    articles: Article[];
    description: string;
}