import useFetch from "../../hooks/useFetch";
import type { Article } from "../../types/article";
import { Link } from "react-router-dom";

export default function ArticlePage() {

    const {data, loading, error} = useFetch<{ data: Article[] }>("/articles?populate=*")
    if (loading) return <main>Loading...</main>;
    if (error) return <main>Error: {error}</main>;

  return (
    <main>
      <h1>Articles</h1>

      {data?.data.map((article) => (
        <div key={article.documentId}>
          {article.cover?.url && (
            <img
              src={`http://localhost:1337${article.cover.url}`}
              alt={article.title}
              width={200}
            />
          )}

          <h2>{article.title}</h2>
          <p>{article.description}</p>

          <p>Slug: {article.slug}</p>
          <p>Author: {article.author?.name}</p>
          <p>Category: {article.category?.name}</p>
          <Link to={`/article/${article.documentId}`}>
            Read More
          </Link>
        </div>
      ))}
    </main>
    );
}

