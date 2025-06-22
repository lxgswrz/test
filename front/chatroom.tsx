"use client"
import React, { useState, useEffect } from "react";
import './styles.css';

const Chatroom = () => {
    const [comments, setComments] = useState([]);
    const [newComment, setNewComment] = useState("");
    const [newName, setNewName] = useState("");
    const [pagination, setPagination] = useState({
        page: 1,
        total: 0
    });
    
    const fetchComments = async (page = 1) => {
        try {
            const url = `http://localhost:8080/comment/get?page=${page}&size=10`;
            const response = await fetch(url);
            
            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.msg || 'Failed to fetch comments');
            }
            
            const result = await response.json();
            
            if (result.code !== 0) {
                throw new Error(result.msg);
            }
            
            const commentsData = Array.isArray(result.data?.comments) 
                ? result.data.comments 
                : [];
            
            setComments(commentsData);
            setPagination({
                page,
                total: result.data?.total || 0
            });
        } catch (err) {
            console.error('Error fetching chat:', err)
        }
    };

    useEffect(() => {
        fetchComments();
    }, []);

    const handleAddComment = async () => {
        if (!newName.trim()) {
            alert("Please enter your name");
            return;
        }
        if (!newComment.trim()) {
            alert("Please enter something");
            return;
        }

        try {
            const response = await fetch('http://localhost:8080/comment/add', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    name: newName,
                    content: newComment
                })
            });

            const result = await response.json();
            
            if (!response.ok || result.code !== 0) {
                throw new Error(result.msg || 'Failed to add comment');
            }

            fetchComments(pagination.page);
            
            setNewComment("");
            setNewName("");
        } catch (err) {
            console.error('Error adding chat:', err)
        }
    };

    const handleDeleteComment = async (id) => {
        try {
            const response = await fetch(`http://localhost:8080/comment/delete?id=${id}`, {
                method: 'POST'
            });

            const result = await response.json();
            
            if (!response.ok || result.code !== 0) {
                throw new Error(result.msg || 'Failed to delete comment');
            }

            fetchComments(pagination.page);
        } catch (err) {
            console.error('Error deleting chat:', err)
        }
    };

    const handlePageChange = (newPage) => {
        const totalPages = Math.ceil(pagination.total / 10);
            
        if (newPage >= 1 && newPage <= totalPages) {
            fetchComments(newPage);
        }
    };

    const totalPages = Math.max(1, Math.ceil(pagination.total / 10));

    return (
        <div className="chat-app">            
            <div className = "input" >
                    <h3>use name</h3>
                    <form>
                        <input 
                            type = "username"
                            value = {newName}
                            onChange = {(e) => setNewName(e.target.value)}
                            placeholder="请输入" 
                            required
                        />
                    </form>
                    <h3>comment content</h3>
                    <form>
                        <input 
                            type = "content"
                            value = {newComment}
                            onChange = {(e) => setNewComment(e.target.value)}
                            placeholder="请输入" 
                            required
                        />
                    </form>
                    <div className = "container">
                        <button type="submit" onClick = {handleAddComment}>提交</button>
                    </div>
                </div>

            <div className="comment-container">
                {comments?.map((comment) => (
                    <div key={comment.id} className="part">
                        <h3>{comment.name}</h3>
                        <p>{comment.content}</p>
                        <div className = "container">
                            <button onClick = {() => handleDeleteComment(comment.id)}> 删除 </button>
                        </div>
                    </div>
                ))}
            </div>
            
            <div className="pagination">
                <button 
                    disabled={pagination.page <= 1}
                    onClick={() => handlePageChange(pagination.page - 1)}
                >
                    前一页
                </button>
                
                <span className="page">Page {pagination.page} of {totalPages}</span>
                
                <button 
                    disabled={pagination.page >= totalPages}
                    onClick={() => handlePageChange(pagination.page + 1)}
                >
                    后一页
                </button>
            </div>
        </div>
    );
};

export default Chatroom;