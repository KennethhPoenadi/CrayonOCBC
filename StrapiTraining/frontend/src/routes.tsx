import { Routes, Route } from "react-router-dom";
import ArticlePage  from "./pages/article/index";
import AuthorPage from "./pages/author";
import CategoryPage from "./pages/category";
import UserPage from "./pages/user";
import ArticleDetailPage from "./pages/article/detail";
import CategoryDetailPage from "./pages/category/detail";
import AuthorDetailPage from "./pages/author/detail";
import UserDetailPage from "./pages/user/detail";
import ArticleCreatePage from "./pages/article/create";


export default function AppRoutes() {
    return ( 
        <Routes>
            <Route path="/article" element={<ArticlePage/>} />
            <Route path="/article/:id" element={<ArticleDetailPage/>} />

            <Route path="/category" element={<CategoryPage/>} />
            <Route path="/category/:id" element={<CategoryDetailPage />} />

            <Route path="/author" element={<AuthorPage/>} />
            <Route path="/author/:id" element={<AuthorDetailPage />} />
            
            <Route path="/user" element={<UserPage/>} />
            <Route path="/user/:id" element={<UserDetailPage/>} />

            <Route path="/article/create" element={<ArticleCreatePage />} />
        </Routes>
    );
}