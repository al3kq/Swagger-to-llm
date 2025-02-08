import requests
import json
from datetime import datetime

API_KEY = "nnqQtN2bOqamegeGq30ogXPtvi8oS7DfYkGHGyjS"
BASE_URL = "https://api.congress.gov/v3"

headers = {'accept': 'application/json'}

def get_current_congress():
    url = f"{BASE_URL}/congress/current?api_key={API_KEY}"
    response = requests.get(url, headers=headers)
    return response.json()['congress']['number']

def get_bill_details(congress, bill_type, bill_number):
    """Fetches detailed information for a specific bill."""
    bill_url = f"{BASE_URL}/bill/{congress}/{bill_type}/{bill_number}?api_key={API_KEY}"
    bill_response = requests.get(bill_url, headers=headers)
    return bill_response.json().get('bill', None) # Return None if 'bill' key is missing

def get_bill_summaries(congress, bill_type, bill_number):
    """Fetches summaries for a specific bill from the API."""
    summaries_url = f"{BASE_URL}/bill/{congress}/{bill_type}/{bill_number}/summaries?api_key={API_KEY}"
    summaries_response = requests.get(summaries_url, headers=headers)
    summaries_data = summaries_response.json()
    if 'summaries' in summaries_data and summaries_data['summaries']:
        return summaries_data['summaries']
    return None

def display_bill_summary(bill_details, summaries):
    """Displays a formatted summary of a bill, including summaries if available."""
    print(f"\n--- {bill_details['title']} ({bill_details['billNumber']}) ---")
    print(f"  Congress: {bill_details['congress']}th")
    print(f"  Bill Type: {bill_details['billType']}")
    print(f"  Introduced: {bill_details['introducedDate']}")
    print(f"  Latest Action: {bill_details['latestAction']['actionDate']} - {bill_details['latestAction']['text']}")

    if summaries:
        print("\n  Summaries:")
        for summary in summaries:
            # Truncate summary for brevity in output
            truncated_summary = summary['text'][:400] + "..." if len(summary['text']) > 400 else summary['text']
            print(f"  - {summary['updateDate']}: {truncated_summary}")
    else:
        print("  No summaries available for this bill.")
    if 'subjects' in bill_details and bill_details['subjects']['legislativeSubjects']:
        subjects = ", ".join([s['name'] for s in bill_details['subjects']['legislativeSubjects']])
        print(f"\n  Subjects: {subjects}")
    else:
        print("  No subjects listed for this bill.")

def explore_environmental_legislation():
    print("\n\n=== Exploring Environmental Legislation Across Time ===")

    # --- Historical Cornerstone: National Park Service Organic Act (64th Congress) ---
    nps_bill_details = get_bill_details(64, 'hr', 15522) # HR15522 in 64th Congress
    if nps_bill_details:
        nps_summaries = get_bill_summaries(64, 'hr', 15522)
        print("\n--- Historical Cornerstone: National Park Service Organic Act of 1916 ---")
        display_bill_summary(nps_bill_details, nps_summaries)
    else:
        print("\nCould not retrieve details for the National Park Service Organic Act.")

    # --- Modern Environmental Bills (Current Congress) ---
    current_congress = get_current_congress()
    environmental_query = "environmental conservation" # Broader search term
    modern_environmental_bills_url = f"{BASE_URL}/summaries?query={environmental_query}&sort=updateDate+desc&limit=3&api_key={API_KEY}"
    modern_bills_response = requests.get(modern_environmental_bills_url, headers=headers).json()

    if 'summaries' in modern_bills_response and modern_bills_response['summaries']:
        print(f"\n\n--- Recent Legislation related to '{environmental_query}' (Current Congress) ---")
        for modern_bill_summary_data in modern_bills_response['summaries']:
            modern_bill_details = get_bill_details(modern_bill_summary_data['congress'], modern_bill_summary_data['billType'], modern_bill_summary_data['billNumber'])
            if modern_bill_details: # Ensure we have bill details before proceeding
                summaries = get_bill_summaries(modern_bill_summary_data['congress'], modern_bill_summary_data['billType'], modern_bill_summary_data['billNumber']) # Fetch summaries again for modern bills

                display_bill_summary(modern_bill_details, summaries)
            else:
                print(f"Could not retrieve full details for {modern_bill_summary_data['title']} ({modern_bill_summary_data['congress']}th Congress)")

    else:
        print(f"\nNo recent summaries found related to '{environmental_query}'.")


def find_ancient_laws():
    # Get laws from the first Congress (1789-1791)
    print("\n\n=== Laws from the 1st Congress (1789-1791) ===")
    first_congress_url = f"{BASE_URL}/law/1?limit=3&api_key={API_KEY}" # Reduced limit for brevity
    response = requests.get(first_congress_url, headers=headers)

    if 'bills' in response.json():
        for law in response.json()['bills']:
            print(f"\n--- Law: {law['title']} ---")
            print(f"  Enacted: {law['enactedDate']}")
            if 'subjects' in law and law['subjects']['legislativeSubjects']:
                subjects = ", ".join([s['name'] for s in law['subjects']['legislativeSubjects']])
                print(f"  Subjects: {subjects}")
            else:
                print("  No subjects listed for this law.")

def main():
    print("Exploring Congressional Time Capsule...\n")
    explore_environmental_legislation() # Focus on environmental theme
    find_ancient_laws()

if __name__ == "__main__":
    main()