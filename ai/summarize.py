import json
import sys

# from transformers import T5Tokenizer, T5ForConditionalGeneration


def summarize(text, model_name="t5-small"):
    """
    Generate summary of the given text using the specified model.

    Args:
    text (str): The text to summarize.
    model_name (str): The name of the model to use for summarization.

    Returns:
    str: The summary of the given text.
    """
    # tokenizer = T5Tokenizer.from_pretrained(model_name)
    # model = T5ForConditionalGeneration.from_pretrained(model_name)
    #
    # input_text = "summarize: " + text
    # input_ids = tokenizer.encode(input_text, return_tensors="pt", max_length=512, truncation=True)
    #
    # summary_ids = model.generate(input_ids, max_length=150, min_length=40, length_penalty=2.0, num_beams=4, early_stopping=True)
    # summary = tokenizer.decode(summary_ids[0], skip_special_tokens=True)
    summary = text
    return summary


if __name__ == "__main__":
    input_text = sys.stdin.read()

    summary = summarize(input_text)

    print(json.dumps({"summary": summary}))