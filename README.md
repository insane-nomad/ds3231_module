# DS3231_MODULE

��������� ��� ������ � �������� �� ������ ���������� ����� ��������� ������� DS3231.

����������� ��������:

- ��������� ������� � NTP-������� � ������ ������� �� DS3231.
- ���������� �������� ���� ����� �� ������ ���������� DS3231.

## ������� ������ ���������

��� ������� ��������� � ����������  ```-command set -mod 1```, ��� 1 - ����� ������ (��������: ���� ������� ����� � ������� ��������� �� �������� ��� ������) ����������:
- ����������� � NTP-�������, ������ ������� ������ ����� � ������������ � �������� ���������� DS3231.
- � �������� ```modules``` ��������� ��������� ���� � ��������� ```mod<����� ������>.txt```, ���� ������������ ���� � ����� ������������� ��� ����������� ������������ ���������� �������� ���� ����� � ppm.

��� ������� ��������� � ����������  ```-command compare -mod 1```, ��� 1 - ����� ������ (��������: ���� ������� ����� � ������� ��������� �� �������� ��� ������) ����������:
- ����������� � NTP-�������, ������ ������� ������ �����
- ������������ ����� � ���������� DS3231
- �������������� ������� �� ������� ����� NTP-�������� � �������� �� ���������� DS3231 (� ��������).
- �������������� ������� �� ������� ����� NTP-�������� � �������� ������������� ������ (� ��������). ��� ���� ����� ��������� ������� ����� ������.
- ��� ���������� ��������� ��������� ������� � (�������� ������) �������, ���:
  -- Time from NTP - ���� � ����� � NTP-�������
  -- Time from file - ���� � ����� ������������� ������ � NTP-��������
  -- Time from module - ���� � ����� � DS3231
  -- Sec from file time - �����, ��������� � ������� �������������
  -- Accuracy - ��������� ���� DS3231 � ppm
   
```	
	+-----------------------------------------------------+
	|                  Time is compared                   |
	+-----------------------------------------------------+
	| Time from NTP         2023-10-09 21:07:13 +0000 UTC |
	| Time from file        2023-10-09 10:00:29 +0000 UTC |
	| Time from module      2023-10-09 21:07:13 +0000 UTC |
	| Sec from file time    40004                     sec |
	| Accuracy              0                         ppm |
	+-----------------------------------------------------+
	
```

## ������������

��� ������ ������� ������� � ������� ������������� �� ������� �������� ����� (�������� Sec from file time) - ��� ���� �������� �������� �������� ���� �����.
������������� �������� Sec from file time �� ����� 1 000 000 ������ (�� ���� �� ����� 12 ���� � ������� ������������� �������), ����� - ������