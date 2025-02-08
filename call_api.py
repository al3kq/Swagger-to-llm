import openai
import os

# --- Configuration Section ---
# 1. Set your OpenAI API Key:
#    It's best to set this as an environment variable for security.
#    You can set it in your terminal like this (replace YOUR_API_KEY):
#    export OPENAI_API_KEY="YOUR_API_KEY"
#    Or directly here (less secure, but for testing):
#    openai.api_key = "YOUR_API_KEY"

# If you are using environment variables (recommended):
openai.api_key = os.getenv("OPENAI_API_KEY")
if not openai.api_key:
    raise ValueError("Please set your OpenAI API key as an environment variable 'OPENAI_API_KEY' or directly in the script.")


# 2. File paths:
input_file_path = "llm_small.txt"  # Path to your input .txt file
output_file_path = "openai_response.txt"  # Path to save the API response

# 3. The string you want to modify and add as the last line of the prompt:
modified_string_line = "use the api doc above to generate a script I will find interesting. I have an API key. Show me find important insights about real things happening recently. Do more than just the summary text. Try to connect things and give me a truly interesting script."

# 4. OpenAI Model to use (you can experiment with different models):
model_name = "o1"  # Or "gpt-4", etc.


# --- Script Logic ---
def create_openai_prompt(txt_file_path, modified_line):
    """Reads the text file and appends the modified line to create the full prompt."""
    try:
        with open(txt_file_path, "r") as file:
            text_content = file.read()
    except FileNotFoundError:
        raise FileNotFoundError(f"Input file not found at: {txt_file_path}")

    full_prompt = text_content.strip() + "\n" + modified_line.strip()
    return full_prompt


def send_prompt_to_openai(prompt_text, model):
    """Sends the prompt to the OpenAI API and returns the response text."""
    try:
        response = openai.ChatCompletion.create(
            model=model,
            messages=[
                {"role": "system", "content": "You are a helpful assistant."},
                {"role": "user", "content": prompt_text},
            ],
            max_completion_tokens=10000,  # Adjust as needed
            n=1,             # Number of completions to generate
            stop=None,       # Stop sequences (optional)
        )
        return response.choices[0].message['content'].strip()
    except openai.error.OpenAIError as e:
        print(f"Error communicating with OpenAI API: {e}")
        return None


def save_response_to_file(response_text, output_path):
    """Saves the API response to the specified output file."""
    if response_text:
        try:
            with open(output_path, "w") as file:
                file.write(response_text)
            print(f"API response saved to: {output_path}")
        except Exception as e:
            print(f"Error saving response to file: {e}")
    else:
        print("No response text to save.")


if __name__ == "__main__":
    try:
        full_prompt = create_openai_prompt(input_file_path, modified_string_line)
        if full_prompt:
            print("--- Prompt sent to OpenAI ---")
            print(full_prompt)  # Optional: print the prompt you're sending

            api_response = send_prompt_to_openai(full_prompt, model_name)

            if api_response:
                print("\n--- OpenAI Response ---")
                print(api_response) # Optional: print the response in the console
                save_response_to_file(api_response, output_file_path)
            else:
                print("No response received from OpenAI.")
        else:
            print("Could not create a valid prompt.")

    except Exception as e:
        print(f"An error occurred: {e}")