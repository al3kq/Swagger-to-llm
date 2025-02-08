import sys
import argparse
from transformers import AutoTokenizer

"""
This script compares the token counts of two files, from the perspective of "file1".

- If file2 has more tokens than file1, it will report how much larger file2 is in percentage.
- If file2 has fewer tokens, it will report how much smaller file2 is.

"""
def count_tokens(file_path, tokenizer):
    try:
        with open(file_path, 'r', encoding='utf-8') as file:
            text = file.read()
            tokens = tokenizer(text)["input_ids"]
            return len(tokens)
    except Exception as e:
        print(f"Error reading file {file_path}: {e}")
        return None


def compare_files(file1, file2, model_name="gpt2"):
    tokenizer = AutoTokenizer.from_pretrained(model_name)
    tokens1 = count_tokens(file1, tokenizer)
    tokens2 = count_tokens(file2, tokenizer)

    if tokens1 is None or tokens2 is None:
        print("Error processing one or both files.")
        return

    print(f"Token count in {file1}: {tokens1}")
    print(f"Token count in {file2}: {tokens2}")

    if tokens1 == 0:
        print("Cannot compute percentage change because file1 has 0 tokens.")
        return

    # From the perspective of file1 => (tokens2 / tokens1 - 1) * 100
    percentage_change = ((tokens2 / tokens1) - 1) * 100

    if percentage_change > 0:
        print(f"{file2} is {percentage_change:.2f}% larger than {file1}.")
    elif percentage_change < 0:
        print(f"{file2} is {abs(percentage_change):.2f}% smaller than {file1}.")
    else:
        print(f"Both files have the same token count.")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Compare the number of tokens in two files.")
    parser.add_argument("file1", help="Path to the first file.")
    parser.add_argument("file2", help="Path to the second file.")
    parser.add_argument("--model", default="gpt2", help="Model name for tokenization (default: gpt2).")
    args = parser.parse_args()

    compare_files(args.file1, args.file2, args.model)
