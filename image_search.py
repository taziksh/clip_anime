import sentence_transformers
from sentence_transformers import SentenceTransformer, util
from PIL import Image
import glob
import torch
import os
from tqdm import tqdm



model = SentenceTransformer('clip-ViT-B-32')

use_pre_computed_embeddings = False
query="Green man"
k=3

if not use_pre_computed_embeddings:
    img_names = list(glob.glob('./*.jpg'))
    print("images: ", len(img_names))
    img_emb = model.encode([Image.open(filepath) for filepath in img_names], batch_size=128, convert_to_tensor=True, show_progress_bar=True)

        # First, we encode the query (which can either be an image or a text string)
    query_emb = model.encode([query], convert_to_tensor=True, show_progress_bar=False)
    
    # Then, we use the util.semantic_search function, which computes the cosine-similarity
    # between the query embedding and all image embeddings.
    # It then returns the top_k highest ranked images, which we output
    hits = util.semantic_search(query_emb, img_emb, top_k=k)[0]
    
    print("Query:")
    print(query)
    for hit in hits:
        print(img_names[hit['corpus_id']])
