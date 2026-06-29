import type { Author } from "../../types/author"
import useFetch from "../../hooks/useFetch";
import { Link } from "react-router-dom";

export default function AuthorPage() {

        const {data, loading, error} = useFetch<{ data: Author[] }>("/authors?populate=*")
    
        if (loading) return <main>Loading...</main>;
        if (error) return <main>Error: {error}</main>;


    return (
        <main>
            <h1>Author</h1>
            {data?.data.map((author) => (
                <div key={author.id}>
                    {author.avatar?.url && (
                        <img src = {`http://localhost:1337${author.avatar.url}`} alt={author.name} width={200}/>
                    )}
                    <h2>
                        {author.name}
                    </h2>
                    <p>
                        {author.email}
                    </p>
                    <p>
                        <Link to={`/author/${author.documentId}`}>Read More</Link>
                    </p>
                </div>
            ))}
        </main>
  
    );
}