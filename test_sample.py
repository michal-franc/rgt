"""
Sample Python test file for testing rgt with Python support.

To test:
1. Run: rgt --test-type python
2. Edit this file and save it
3. rgt should automatically run pytest on this file
"""

def add(a, b):
    """Simple function to add two numbers."""
    return a + b


def test_add_positive_numbers():
    """Test adding positive numbers."""
    assert add(2, 3) == 5
    assert add(10, 20) == 30


def test_add_negative_numbers():
    """Test adding negative numbers."""
    assert add(-5, -3) == -8
    assert add(-10, 5) == -5


def test_add_zero():
    """Test adding with zero."""
    assert add(0, 0) == 0
    assert add(5, 0) == 5
    assert add(0, 10) == 10


if __name__ == "__main__":
    # Run tests manually if executed directly
    test_add_positive_numbers()
    test_add_negative_numbers()
    test_add_zero()
    print("All tests passed!")
